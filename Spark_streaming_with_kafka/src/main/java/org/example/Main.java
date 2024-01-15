package org.example;

import org.apache.log4j.Level;
import org.apache.log4j.Logger;
import org.apache.spark.SparkConf;
import org.apache.spark.api.java.JavaDoubleRDD;
import org.apache.spark.api.java.function.ForeachFunction;
import org.apache.spark.api.java.function.VoidFunction2;
import org.apache.spark.internal.config.R;
import org.apache.spark.sql.*;
import org.apache.spark.sql.streaming.OutputMode;
import org.apache.spark.sql.streaming.StreamingQuery;
import org.apache.spark.sql.streaming.StreamingQueryException;
import org.apache.spark.sql.types.DataTypes;
import org.apache.spark.sql.types.StructField;
import org.apache.spark.sql.types.StructType;
import org.apache.spark.streaming.Duration;
import org.apache.spark.streaming.Durations;
import org.apache.spark.streaming.api.java.JavaStreamingContext;
import scala.Function2;
import scala.runtime.BoxedUnit;

import javax.xml.crypto.Data;
import java.sql.*;
import java.util.Properties;
import java.util.concurrent.TimeoutException;

import static org.apache.spark.sql.functions.*;
// max congestion rate 12 bikes / min
//


public class Main {
    public static void main(String[] args) throws TimeoutException, StreamingQueryException {
        Logger.getLogger("org.apache").setLevel(Level.WARN);
        Logger.getLogger("org.apache.spark.storage").setLevel(Level.ERROR);
        String url = "jdbc:mysql://localhost:3306/iot_project";
        String username = "guest";
        String password = "2107";
//
        SparkSession session = SparkSession.builder()
                .master("local[*]")
                .appName("iot_streaming")
                .getOrCreate();

        Dataset<Row> df = session.readStream()
                .format("kafka")
                .option("kafka.bootstrap.servers", "localhost:9092")
                .option("subscribe", "quickstart-events")
                .load();
        StructType schema = DataTypes.createStructType(new StructField[] {
                DataTypes.createStructField("parking_lot_id",  DataTypes.StringType, false),
                DataTypes.createStructField("license_plate", DataTypes.StringType, false),
                DataTypes.createStructField("entrance", DataTypes.StringType, false),
        });
        df = df.selectExpr("CAST(value AS STRING)");
        df = df.select(functions.from_json(functions.col("value"), schema).alias("data"));
        df = df.select("data.*");

//        df.createOrReplaceTempView("streaming_data");

//        Dataset<Row> entrance_true = session.sql("select * from streaming_data where entrance = 'True'");
//        Dataset<Row> entrance_false = session.sql("select * from streaming_data where entrance != 'True'");

//        System.out.println();
        df = df.withColumn("timestamp", functions.current_timestamp());

//        Dataset<Row> finalDf = df.withWatermark("timestamp", "5 minutes").groupBy(functions.window(df.col("timestamp"), "10 minutes", "5 minutes"), df.col("parking_lot_id")).count().as("congestion_rate");
//        df.withColumn("count", finalDf.col("count"));
        //        StreamingQuery query = df.writeStream().format("console").option("truncate", "False")
//                .outputMode(OutputMode.Complete()).option("compression","none").start();
//        query.awaitTermination();
        StreamingQuery q = df.writeStream()
//                .outputMode("Complete")
                .foreach(
                        new ForeachWriter<Row>() {
                            @Override
                            public boolean open(long partitionId, long epochId) {

                                return true;
                            }
                            @Override
                            public void process(Row value) {
                                Connection connection = DbConnection.getConnection(url, username, password);

//                                int congestion = Integer.parseInt(value.getAs("count"));
                                int parking_id = Integer.parseInt(value.getAs("parking_lot_id"));
                                boolean is_entering = Boolean.parseBoolean(value.getAs("entrance"));
                                String license_plate =value.getAs("license_plate");

                                System.out.println("row info: " + parking_id + " - " + license_plate + " - " + is_entering);

//                                System.out.println("congestion rate: " + congestion);
                                PreparedStatement preparedStatement = null;

//                                if(congestion > 0){
//                                    String query = "UPDATE parking_lots SET congestionRate = ? WHERE id = ?";
//                                    try {
//                                        preparedStatement = connection.prepareStatement(query);
//                                        preparedStatement.setInt(1, congestion);
//                                        preparedStatement.setInt(2, parking_id);
//                                        preparedStatement.execute();
//                                    } catch (SQLException e) {
//                                        throw new RuntimeException(e);
//                                    }
//                                }

                                if(is_entering){
                                    String query = "INSERT IGNORE INTO license_plates(name, parkinglot_id) VALUES(?, ?)";
                                    try {
                                        preparedStatement = connection.prepareStatement(query);
                                        preparedStatement.setString(1, license_plate);
                                        preparedStatement.setInt(2, parking_id);
                                        System.out.println(preparedStatement.execute());

                                    } catch (SQLException e) {
                                        throw new RuntimeException(e);
                                    }
                                    try {
                                        connection.close();
                                    } catch (SQLException e) {
                                        throw new RuntimeException(e);
                                    }
                                }else{
                                    String query = "DELETE FROM license_plates WHERE name=?";
                                    try {
                                        preparedStatement = connection.prepareStatement(query);
                                        preparedStatement.setString(1, license_plate);
                                        preparedStatement.execute();

                                    } catch (SQLException e) {
                                        throw new RuntimeException(e);
                                    }
                                    try {
                                        connection.close();
                                    } catch (SQLException e) {
                                        throw new RuntimeException(e);
                                    }
                                }
                            }
                            @Override
                            public void close(Throwable errorOrNull) {
                            }
                        }
                ).start();
                q.awaitTermination();
    }
}
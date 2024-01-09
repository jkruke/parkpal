package org.example;

import org.apache.log4j.Level;
import org.apache.log4j.Logger;
import org.apache.spark.sql.*;
import org.apache.spark.sql.streaming.OutputMode;
import org.apache.spark.sql.streaming.StreamingQuery;
import org.apache.spark.sql.streaming.StreamingQueryException;
import org.apache.spark.sql.types.DataTypes;
import org.apache.spark.sql.types.StructField;
import org.apache.spark.sql.types.StructType;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.util.concurrent.TimeoutException;

public class CongestionRate {
    public static void main(String[] args) throws StreamingQueryException, TimeoutException {
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

        df = df.withWatermark("timestamp", "5 minutes").groupBy(functions.window(df.col("timestamp"), "10 minutes", "5 minutes"), df.col("parking_lot_id")).count().as("congestion_rate");
        StreamingQuery query = df.writeStream()
                .outputMode("Complete")
                .foreach(
                new ForeachWriter<Row>() {
                    @Override
                    public boolean open(long partitionId, long epochId) {
                        return true;
                    }

                    @Override
                    public void process(Row value) {
                        Connection connection = DbConnection.getConnection(url, username, password);

                        int congestion = (int) value.getLong(2);
                        int parking_id = Integer.parseInt(value.getAs("parking_lot_id"));
//                        boolean is_entering = Boolean.parseBoolean(value.getAs("entrance"));
//                        String license_plate =value.getAs("license_plate");

//                        System.out.println("row info: " + parking_id + " - " + license_plate + " - " + is_entering);

                                System.out.println("congestion rate: " + congestion);
                        PreparedStatement preparedStatement = null;

                                if(congestion > 0){
                                    String query = "UPDATE parking_lots SET congestionRate = ? WHERE id = ?";
                                    try {
                                        preparedStatement = connection.prepareStatement(query);
                                        preparedStatement.setInt(1, congestion);
                                        preparedStatement.setInt(2, parking_id);
                                        preparedStatement.execute();
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
        query.awaitTermination();
    }
}

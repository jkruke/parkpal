package org.example;

import java.sql.*;

public class Query_sqlite {
    private static boolean execute;

    public static void main(String[] argv) {
        // SQLite connection string
        String url = "jdbc:sqlite:/home/kytran/Downloads/Spark_streaming_with_kafka/database.db";

        // SQL statement for creating a new tablex`
        String sql = "select * from license_plate_tbl";

        try (Connection conn = DriverManager.getConnection(url);
             Statement stmt = conn.createStatement()) {
            // create a new table
            ResultSet rs = stmt.executeQuery(sql);
            while (rs.next()) {
                System.out.println(rs.getInt("id") +  "\t" +
                        rs.getString("license_plate") + "\t" +
                        rs.getDate("timestamp"));
            }

        } catch (SQLException e) {
            System.out.println(e.getMessage());
        }
    }
}

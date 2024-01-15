package org.example;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;

public class DbConnection {
    public static Connection getConnection(String url, String username, String password) {
        Connection connection = null;
        try {
//            String url = "jdbc:mysql://localhost:3306/iot_project";
//            String username = "guest";
//            String password = "2107";
            connection = DriverManager.getConnection(url, username, password);
            System.out.println("Connected to the database!");
        } catch (SQLException e) {
            System.err.println("Connection error: " + e.getMessage());
        }
        return connection;
    }

//    public static Connection getConnection(String dbUrl, String userName, String password) {
//    }
}

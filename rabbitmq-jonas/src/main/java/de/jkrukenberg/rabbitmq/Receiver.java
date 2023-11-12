package de.jkrukenberg.rabbitmq;

import com.rabbitmq.client.Channel;
import com.rabbitmq.client.DeliverCallback;

import java.io.IOException;
import java.nio.charset.StandardCharsets;

public class Receiver {

    public static void main(String[] args) {
        AMQPChannelHandler.operate(Receiver::handleChannel);
    }

    private static void handleChannel(Channel channel) {
        DeliverCallback deliverCallback = (consumerTag, delivery) -> {
            String message = new String(delivery.getBody(), StandardCharsets.UTF_8);
            System.out.println(" [x] Received '" + message + "'");
        };
        try {
            channel.basicConsume(AMQPChannelHandler.QUEUE_NAME, true, deliverCallback, consumerTag -> {
            });
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        System.out.println(" [*] Waiting for messages. To exit press CTRL+C");
    }
}

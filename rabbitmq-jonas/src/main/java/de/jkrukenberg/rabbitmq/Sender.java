package de.jkrukenberg.rabbitmq;

import com.rabbitmq.client.Channel;

import java.io.IOException;

public class Sender {

    public static void main(String[] args) {
        AMQPChannelHandler.operate(Sender::handleChannel);
    }

    private static void handleChannel(Channel channel) {
        try {
            String message = "{\"id\": 1, \"plate_number\", \"ABC-123\", \"owner\": \"me\"}";
            channel.basicPublish("", AMQPChannelHandler.QUEUE_NAME, null, message.getBytes());
            System.out.println(" [x] Sent '" + message + "'");
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}

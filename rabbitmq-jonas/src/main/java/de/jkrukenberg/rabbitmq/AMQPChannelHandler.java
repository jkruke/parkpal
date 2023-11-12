package de.jkrukenberg.rabbitmq;

import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;

import java.io.IOException;
import java.util.concurrent.TimeoutException;
import java.util.function.Consumer;

class AMQPChannelHandler {

    static final String QUEUE_NAME = "license-plate-detections";
    private static final String HOST = "localhost";
    private static final int PORT = ConnectionFactory.DEFAULT_AMQP_PORT;

    static void operate(Consumer<Channel> channelConsumer) {
        ConnectionFactory factory = new ConnectionFactory();
        factory.setHost(HOST);
        factory.setPort(PORT);
        try (Connection connection = factory.newConnection()) {
            Channel channel = connection.createChannel();
            channel.queueDeclare(QUEUE_NAME, false, false, false, null);

            channelConsumer.accept(channel);

        } catch (IOException | TimeoutException e) {
            throw new RuntimeException(e);
        }
    }
}

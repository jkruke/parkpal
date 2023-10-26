import org.eclipse.paho.client.mqttv3.*;
import org.eclipse.paho.client.mqttv3.persist.MemoryPersistence;

public class App {
    
    private static final String MQTT_BROKER = "tcp://broker.hivemq.com:1883";
    private static final String MQTT_SUB_TOPIC = "/hust-iot/jonas";
    private static final String MQTT_PUB_TOPIC = MQTT_SUB_TOPIC;
    private static final String MQTT_CONTENT = "{\"id\": 1, \"plate_number\", \"ABC-123\", \"owner\": \"me\"}";
    private static final int MQTT_QOS = 1;

    public static void main(String[] args) {
        new App().start();
    }

    public void start() {
        String clientId = "ID_OF_CLIENT";
        MemoryPersistence persistence = new MemoryPersistence();
        try {
            MqttClient client = new MqttClient(MQTT_BROKER, clientId, persistence);
            MqttConnectOptions connOpts = new MqttConnectOptions();
            connOpts.setCleanSession(true);
            client.connect(connOpts);

            client.setCallback(getMqttCallback());
            client.subscribe(MQTT_SUB_TOPIC);

            MqttMessage message = new MqttMessage(MQTT_CONTENT.getBytes());
            message.setQos(MQTT_QOS);
            System.out.println("Publishing message: " + MQTT_CONTENT);
            client.publish(MQTT_PUB_TOPIC, message);
            System.out.println("Message published");
            client.disconnect();
            System.out.println("Disconnected");
            client.close();
        } catch (MqttException e) {
            throw new RuntimeException(e);
        }
    }

    private static MqttCallback getMqttCallback() {
        return new MqttCallback() {
            @Override
            public void connectionLost(Throwable throwable) {
                System.out.println("Disconnected，you can reconnect.");
            }

            @Override
            public void messageArrived(String s, MqttMessage mqttMessage) {
                System.out.println("Received message: " + new String(mqttMessage.getPayload()));
            }

            @Override
            public void deliveryComplete(IMqttDeliveryToken iMqttDeliveryToken) {
                System.out.println("Delivery completed.");
            }
        };
    }
}

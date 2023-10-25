import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;

import java.io.IOException;
import java.io.OutputStream;

public class SimpleHandler implements HttpHandler {

    @Override
    public void handle(HttpExchange he) throws IOException {
        String response = "{\"id\": 1, \"plate_number\", \"ABC-123\", \"owner\": \"me\"}";
        he.sendResponseHeaders(200, response.length());
        OutputStream os = he.getResponseBody();
        os.write(response.getBytes());
    }
}
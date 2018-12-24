/**
 * @author Fabio Colacio
 */

package xyz.fabiocolacio.quicksilver;

import java.net.URL;
import java.net.MalformedURLException;
import java.net.URLEncoder;
import java.io.IOException;
import java.io.OutputStream;
import java.io.OutputStreamWriter;
import java.io.InputStreamReader;
import javax.net.ssl.HttpsURLConnection;
import java.security.Key;
import java.security.NoSuchAlgorithmException;
import java.security.InvalidKeyException;
import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import org.bouncycastle.crypto.generators.SCrypt;
import org.json.JSONWriter;
import org.json.JSONObject;
import org.json.JSONException;
import android.util.Base64;
import android.util.Log;

public class API {
    private static String host = "https://fabiocolacio.xyz:443";

    private static int SCRYPT_CPU_MEM_COST = 32768;
    private static int SCRYPT_BLOCK_SIZE = 8;
    private static int SCRYPT_PARALLELIZATION = 1;
    private static int SALTED_HASH_LEN = 32;

    public static String getHost() {
        return host;
    }

    public static int register(String username, String password) {
        try {
            URL url = new URL(host + "/register");

            HttpsURLConnection con = (HttpsURLConnection) url.openConnection();
            con.setRequestMethod("POST");
            con.setDoOutput(true);

            OutputStreamWriter writer = new OutputStreamWriter(con.getOutputStream());

            new JSONWriter(writer).object()
                .key("Username").value(username)
                .key("Password").value(password)
            .endObject();

            writer.flush();
            writer.close();

            int status = con.getResponseCode();
            return status;
        } catch (MalformedURLException e) {
            System.out.println(e);
        } catch (JSONException e) {
            System.out.println(e);
        } catch (IOException e) {
            System.out.println(e);
        }

        return -1;
    }

    public static String login(String username, String password) {
        try {
            String encodedUsername = URLEncoder.encode(username, "UTF-8");
            URL url = new URL(host + "/login?user=" + encodedUsername);

            HttpsURLConnection con = (HttpsURLConnection) url.openConnection();
            con.setRequestMethod("POST");

            int status = con.getResponseCode();
            if (status != 200) {
                Log.e("Login", "auth failed");
                return null;
            }

            int contentLength = con.getContentLength();
            char[] body = new char[contentLength];
            InputStreamReader reader = new InputStreamReader(con.getInputStream());
            reader.read(body, 0, contentLength);
            reader.close();
        
            JSONObject json = new JSONObject(new String(body));
            String challengeString = json.getString("C");
            String saltString = json.getString("S");
            System.out.printf("Challenge: %s\nSalt: %s\n", challengeString, saltString);

            byte[] challenge = Base64.decode(challengeString, Base64.DEFAULT);
            byte[] salt = Base64.decode(saltString, Base64.DEFAULT);

            byte[] saltedHash = SCrypt.generate(
                password.getBytes(),
                salt,
                SCRYPT_CPU_MEM_COST,
                SCRYPT_BLOCK_SIZE,
                SCRYPT_PARALLELIZATION,
                SALTED_HASH_LEN);
       
            Key macKey = new SecretKeySpec(saltedHash, "HmacSHA256");
            Mac mac = Mac.getInstance("HmacSHA256");
            mac.init(macKey);
            byte[] response = mac.doFinal(challenge);

            url = new URL(host + "/auth?user=" + encodedUsername);

            con = (HttpsURLConnection) url.openConnection();
            con.setRequestMethod("POST");
            con.setDoOutput(true);

            OutputStream outputStream = con.getOutputStream();
            outputStream.write(response);
            outputStream.flush();
            outputStream.close();

            status = con.getResponseCode();
            if (status != 200) {
                Log.e("Login", "auth failed");
                return null;
            }

            contentLength = con.getContentLength();
            char[] jwt = new char[contentLength];
            reader = new InputStreamReader(con.getInputStream());
            reader.read(body, 0, contentLength);
            reader.close();

            Log.i("Login", "JWT: " + new String(jwt));

            return new String(jwt);
        } catch (MalformedURLException e) {
            System.out.println(e);
        } catch (NoSuchAlgorithmException e) {
            System.out.println(e);
        } catch (InvalidKeyException e) {
            System.out.println(e);
        } catch (JSONException e) {
            System.out.println(e);
        } catch (IOException e) {
            System.out.println(e);
        }

        return null;
    }
}


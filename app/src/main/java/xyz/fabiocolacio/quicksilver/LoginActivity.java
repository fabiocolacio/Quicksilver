package xyz.fabiocolacio.quicksilver;

import android.app.ActionBar;
import android.os.AsyncTask;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.EditText;

import xyz.fabiocolacio.quicksilver.API;

public class LoginActivity extends AppCompatActivity {
    private class LoginTask extends AsyncTask<String, Void, String> {
        @Override
        protected String doInBackground(String... params) {
            if (params.length < 2) {
                return null;
            }

            return API.login(params[0], params[1]);
        }

        @Override
        protected void onPostExecute(String jwt) {
            Log.i("Login", "JWT: " + jwt);
        }
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_login);

        ActionBar bar = getActionBar();
        if (bar != null) {
            bar.setTitle("Login");
        }


        android.support.v7.app.ActionBar supportActionBar = getSupportActionBar();
        if (supportActionBar != null) {
            supportActionBar.setTitle("Login");
        }
    }

    public void loginClicked(View button) {
        EditText userEntry = (EditText) findViewById(R.id.userLoginEntry);
        EditText passEntry = (EditText) findViewById(R.id.passLoginEntry);

        String user = userEntry.getText().toString();
        String pass = passEntry.getText().toString();

        new LoginTask().execute(user, pass);
    }
}

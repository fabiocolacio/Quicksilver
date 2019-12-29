package xyz.fabiocolacio.quicksilver;

import android.app.ActionBar;
import android.os.AsyncTask;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.EditText;

public class RegisterActivity extends AppCompatActivity {
    private class RegisterTask extends AsyncTask<String, Void, Integer> {
        @Override
        protected Integer doInBackground(String... params) {
            if (params.length < 2) {
                return null;
            }

            return API.register(params[0], params[1]);
        }

        @Override
        protected void onPostExecute(Integer status) {
            Log.i("Register", (status == 200) ? "Registration successful" : "Registration failed");
        }
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_register);

        ActionBar bar = getActionBar();
        if (bar != null) {
            bar.setTitle("Register");
        }


        android.support.v7.app.ActionBar supportActionBar = getSupportActionBar();
        if (supportActionBar != null) {
            supportActionBar.setTitle("Register");
        }
    }

    public void registerClicked(View button) {
        EditText userEntry = (EditText) findViewById(R.id.userRegisterEntry);
        EditText passEntry = (EditText) findViewById(R.id.passRegisterEntry);

        String user = userEntry.getText().toString();
        String pass = passEntry.getText().toString();

        new RegisterTask().execute(user, pass);
    }
}

package xyz.fabiocolacio.quicksilver;

import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;
import xyz.fabiocolacio.quicksilver.API;

public class Database extends SQLiteOpenHelper {
    private static final String DATABASE_NAME = "quicksilver";
    private static final int DATABASE_VERSION = 1;

    public enum Table {
        PUBKEYS("privkeys",
                "CREATE TABLE privkeys(" +
                        "id INT," +
                        "owner VARCHAR(" + API.USERNAME_MAX_LEN + ")," +
                        "peer VARCHAR(" + API.USERNAME_MAX_LEN + ")," +
                        "privkey BLOB," +
                        "PRIMARY KEY (id, owner, peer))"),

        PRIVKEYS("privkeys",
                 "CREATE TABLE pubkeys(" +
                        "id INT," +
                        "owner VARCHAR(" + API.USERNAME_MAX_LEN + ")," +
                        "peer VARCHAR(" + API.USERNAME_MAX_LEN + ")," +
                        "pubkey BLOB," +
                        "PRIMARY KEY (id, owner, peer))");

        private String name;
        private String def;

        Table(String name, String def) {
            this.name = name;
            this.def = def;
        }

        public String getName() {
            return this.name;
        }

        public String getDef() {
            return this.def;
        }
    };

    Database(Context context) {
        super(context, DATABASE_NAME, null, DATABASE_VERSION);
    }

    @Override
    public void onCreate(SQLiteDatabase db) {
        for (Table table : Table.values()) {
            db.execSQL(table.getDef());
        }
    }

    @Override
    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {
        for (Table table: Table.values()) {
            db.execSQL("DROP TABLE IF EXISTS " + table.getName());
        }
        onCreate(db);
    }
}

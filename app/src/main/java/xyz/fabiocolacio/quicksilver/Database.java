package xyz.fabiocolacio.quicksilver;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
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

    public byte[] lookupPubKey(String owner, String peer, int id) {
        SQLiteDatabase db = this.getReadableDatabase();

        String[] args = { owner, peer, String.valueOf(id) };

        Cursor cursor = db.rawQuery(
                "SELECT pubkey" +
                    "FROM pubkeys" +
                    "WHERE owner = ? AND peer = ? AND id = ?",
                    args);

        byte[] pubKey = cursor.getBlob(0);

        cursor.close();
        db.close();

        return pubKey;
    }

    public byte[] lookupPrivKey(String owner, String peer, int id) {
        SQLiteDatabase db = this.getReadableDatabase();

        String[] args = { owner, peer, String.valueOf(id) };

        Cursor cursor = db.rawQuery(
                "SELECT privkey" +
                    "FROM privkeys" +
                    "WHERE owner = ? AND peer = ? AND id = ?",
                args);

        byte[] privKey = cursor.getBlob(0);

        cursor.close();
        db.close();

        return privKey;
    }

    public void uploadPubKey(String owner, String peer, int id, byte[] pubKey) {
        SQLiteDatabase db = this.getWritableDatabase();

        ContentValues values = new ContentValues();
        values.put("pubkey", pubKey);
        values.put("id", id);
        values.put("owner", owner);
        values.put("peer", peer);

        db.insert("pubkeys", null, values);

        db.close();
    }

    public void uploadPrivKey(String owner, String peer, int id, byte[] privKey) {
        SQLiteDatabase db = this.getWritableDatabase();

        ContentValues values = new ContentValues();
        values.put("privkey", privKey);
        values.put("id", id);
        values.put("owner", owner);
        values.put("peer", peer);

        db.insert("privkeys", null, values);

        db.close();
    }

    public int latestPubKey(String owner, String peer) {
        SQLiteDatabase db = this.getReadableDatabase();

        String[] args = { owner, peer };

        Cursor cursor = db.rawQuery(
                "SELECT id" +
                        "FROM pubkeys" +
                        "WHERE owner = ? AND peer = ?" +
                        "ORDER BY id DESC",
                args);

        int id = cursor.getInt(0);

        cursor.close();
        db.close();

        return id;
    }

    public int latestPrivKey(String owner, String peer) {
        SQLiteDatabase db = this.getReadableDatabase();

        String[] args = { owner, peer };

        Cursor cursor = db.rawQuery(
                "SELECT id" +
                    "FROM privkeys" +
                    "WHERE owner = ? AND peer = ?" +
                    "ORDER BY id DESC",
                    args);

        int id = cursor.getInt(0);

        cursor.close();
        db.close();

        return id;
    }
}

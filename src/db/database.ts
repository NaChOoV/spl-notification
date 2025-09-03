import { drizzle } from 'drizzle-orm/libsql';

import { migrate } from 'drizzle-orm/libsql/migrator';

import EnvConfig from '../config/enviroment';

const db = drizzle({
    connection: {
        url: EnvConfig.databaseUrl,
        authToken: EnvConfig.databaseToken,
    },
});

export const runMigration = async () => {
    await migrate(db, { migrationsFolder: './drizzle' });
    console.log('[Drizzle] Success Database Migration...');
};

export default db;

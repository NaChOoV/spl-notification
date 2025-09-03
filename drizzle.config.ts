import EnvConfig from './src/config/enviroment';
import { defineConfig } from 'drizzle-kit';

export default defineConfig({
    out: './drizzle',
    schema: './src/db/schema.ts',
    dialect: 'turso',
    dbCredentials: {
        url: EnvConfig.databaseUrl,
        authToken: EnvConfig.databaseToken,
    },
});

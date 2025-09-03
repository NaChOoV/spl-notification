const EnvConfig = {
    port: process.env.PORT || 4001,
    authString: process.env.AUTH_STRING || '',
    timeZone: process.env.TIME_ZONE || 'GMT-3',
    databaseUrl: process.env.TURSO_DATABASE_URL || '',
    databaseToken: process.env.TURSO_AUTH_TOKEN || '',
    whatsappBaseUrl: process.env.WHATSAPP_BASE_URL || '',
    whatsappUsername: process.env.WHATSAPP_USERNAME || '',
    whatsappPassword: process.env.WHATSAPP_PASSWORD || '',
    sourceServiceBaseUrl: process.env.SOURCE_BASE_URL || '',
    sourceServiceAuthString: process.env.SOURCE_AUTH_STRING || '',
    accessServiceBaseUrl: process.env.ACCESS_SERVICE_BASE_URL || '',
    accessServiceAuthString: process.env.ACCESS_SERVICE_AUTH_STRING || '',
};

export default EnvConfig;

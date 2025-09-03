import type { Context, Next } from 'hono';
import EnvConfig from '../config/enviroment';

export const authMiddleware = async (c: Context, next: Next) => {
    const authHeader = c.req.header('X-Auth-Token');
    if (authHeader !== EnvConfig.authString) {
        return c.text('Unauthorized', 401);
    }
    await next();
};

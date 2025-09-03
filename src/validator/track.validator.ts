import { zValidator } from '@hono/zod-validator';
import * as z from 'zod';

const trackValidator = zValidator(
    'json',
    z.object({
        chatId: z.string(),
        run: z.string(),
        alias: z.string().optional(),
    })
);

export { trackValidator };

import { zValidator } from '@hono/zod-validator';
import * as z from 'zod';

const accessValidator = zValidator(
    'json',
    z.array(
        z.object({
            externalId: z.string(),
            run: z.string(),
            fullName: z.string(),
            location: z.string(),
            entryAt: z.string(),
        })
    )
);

export { accessValidator };

import { Hono } from 'hono';
import { serve } from '@hono/node-server';
import EnvConfig from './config/enviroment';
import { authMiddleware } from './middleware/auth';
import { runMigration } from './db/database';
import type { Access } from './types/access';
import { accessService, sourceService, trackRepository, whatsappService } from './di-container';
import { accessValidator } from './validator/access.validator';
import { TrackType } from './db/schema';
import { HttpStatusCode } from 'axios';
import type { TrackDto } from './dto/track.dto';
import { trackValidator } from './validator/track.validator';
import { logger } from 'hono/logger';

await runMigration();

const app = new Hono();
app.use(
    logger((str, ...rest) => {
        console.log(new Date().toISOString(), str, ...rest);
    })
);

app.post('/access', authMiddleware, accessValidator, async (c) => {
    const access = await c.req.json<Access[]>();

    await accessService.checkAccess(access);
    return c.body(null, HttpStatusCode.Created);
});

app.post('/track', authMiddleware, trackValidator, async (c) => {
    const trackDto = await c.req.json<TrackDto>();

    const abmUser = await sourceService.getAbmUserByRun(trackDto.run);
    if (!abmUser) {
        await whatsappService.sendMessage(trackDto.chatId, 'Usuario no existente');
        return c.body(null, HttpStatusCode.NotFound);
    }

    const track = await trackRepository.listTrack(Number(trackDto.chatId), TrackType.TRACK);
    const found = track.some((t) => t.userId === String(abmUser.externalId));
    if (found) {
        await whatsappService.sendMessage(trackDto.chatId, '✅ Seguimiento agregado');
        return c.body(null, HttpStatusCode.Ok);
    }

    const accesses = await accessService.getAccess();
    const access = accesses.find((access) => access.externalId === String(abmUser.externalId));
    const lastEntry = access?.entryAt;

    const newTrack = {
        chatId: Number(trackDto.chatId),
        userId: String(abmUser.externalId),
        run: trackDto.run,
        fullName: `${abmUser.firstName} ${abmUser.lastName}`,
        type: TrackType.TRACK,
        alias: trackDto.alias,
        lastEntry,
    };

    await trackRepository.createTrack(newTrack);

    await whatsappService.sendMessage(trackDto.chatId, '✅ Seguimiento agregado');
    return c.body(null, HttpStatusCode.Ok);
});

app.delete('/track', authMiddleware, trackValidator, async (c) => {
    const trackDto = await c.req.json<TrackDto>();
    await trackRepository.removeTrack(Number(trackDto.chatId), trackDto.run, TrackType.TRACK);

    await whatsappService.sendMessage(trackDto.chatId, '✅ Seguimiento eliminado');
    return c.body(null, HttpStatusCode.Ok);
});

app.get('/track/:chatId', authMiddleware, async (c) => {
    const chatId = c.req.param('chatId');

    const tracks = await trackRepository.listTrack(Number(chatId), TrackType.TRACK);
    await whatsappService.listTracks(chatId, tracks);

    return c.body(null, HttpStatusCode.Ok);
});

serve({
    fetch: app.fetch,
    port: Number(EnvConfig.port),
});

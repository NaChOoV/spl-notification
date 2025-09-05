import { and, or, eq, ne, SQL, sql, inArray } from 'drizzle-orm';
import db from '../db/database';
import { track, TrackType, type Track } from '../db/schema';
import type { Access } from '../types/access';

class TrackRepository {
    private readonly client: typeof db;
    constructor() {
        this.client = db;
    }

    public async createTrack(trackData: {
        chatId: number;
        userId: string;
        run: string;
        fullName: string;
        type: TrackType;
        alias?: string;
        lastEntry?: string;
    }): Promise<void> {
        await this.client.insert(track).values(trackData).onConflictDoNothing();
    }

    public listTrack(chatId: number, type: TrackType): Promise<Track[]> {
        return this.client
            .select()
            .from(track)
            .where(and(eq(track.chatId, chatId), eq(track.type, type)));
    }

    public async removeTrack(chatId: number, run: string, type: TrackType): Promise<void> {
        await this.client
            .delete(track)
            .where(
                and(
                    eq(track.chatId, chatId),
                    sql`UPPER(${track.run}) = UPPER(${run})`,
                    eq(track.type, type)
                )
            );
    }

    public async checkTrack(accessToVerify: Access[], type: TrackType): Promise<Track[]> {
        const conditionals = accessToVerify.reduce<SQL[]>((acc, access) => {
            const condition = and(
                eq(track.userId, access.externalId),
                ne(track.lastEntry, access.entryAt),
                eq(track.type, type)
            );
            return condition ? [...acc, condition] : acc;
        }, []);

        const trackResponse = await this.client
            .select()
            .from(track)
            .where(or(...conditionals));

        console.log(trackResponse);
        return trackResponse;
    }

    public async updateTrack(accesses: Access[], type: TrackType): Promise<void> {
        const sqlChunks: SQL[] = [];
        const userIds: string[] = [];

        sqlChunks.push(sql`(case`);
        for (const access of accesses) {
            sqlChunks.push(sql`when ${track.userId} = ${access.externalId} then ${access.entryAt}`);
            userIds.push(access.externalId);
        }
        sqlChunks.push(sql`end)`);

        const finalSql: SQL = sql.join(sqlChunks, sql.raw(' '));

        await db
            .update(track)
            .set({ lastEntry: finalSql })
            .where(and(inArray(track.userId, userIds), eq(track.type, type)));
    }

    public async getAll(type: TrackType): Promise<Track[]> {
        return this.client.select().from(track).where(eq(track.type, type));
    }

    public async getTrackByUserId(userId: string, type: TrackType): Promise<Track[]> {
        return this.client
            .select()
            .from(track)
            .where(and(eq(track.type, type), eq(track.userId, userId)));
    }
}

export default TrackRepository;

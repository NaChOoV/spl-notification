import { int, sqliteTable, text, integer, unique } from 'drizzle-orm/sqlite-core';
import { type InferSelectModel } from 'drizzle-orm';

export enum TrackType {
    TRACK = '1',
    DELETE = '2',
}

export const track = sqliteTable(
    'track',
    {
        id: int().primaryKey({ autoIncrement: true }),
        chatId: integer().notNull(),
        userId: text().notNull(),
        run: text().notNull(),
        fullName: text().notNull(),
        alias: text(),
        lastEntry: text().default(''),
        lastExit: text().default(''),
        type: text().default(TrackType.TRACK).notNull(),
    },
    (t) => [unique().on(t.chatId, t.run, t.type)]
);

export type Track = InferSelectModel<typeof track>;

export type TrackWithLocation = { location: string } & Track;

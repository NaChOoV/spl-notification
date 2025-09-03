CREATE TABLE `track` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`chatId` integer NOT NULL,
	`userId` text NOT NULL,
	`run` text NOT NULL,
	`fullName` text NOT NULL,
	`lastEntry` text DEFAULT '',
	`type` text DEFAULT '1' NOT NULL
);
--> statement-breakpoint
CREATE UNIQUE INDEX `track_chatId_run_type_unique` ON `track` (`chatId`,`run`,`type`);
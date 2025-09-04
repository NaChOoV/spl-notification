import { CronJob } from 'cron';

import { getSleepSeconds } from '../utils/sleep-time';
import EnvConfig from '../config/enviroment';
import { sleep } from 'bun';
import type AccessService from './access.service';

class CronJobService {
    constructor(private readonly accessService: AccessService) {}

    public setup() {
        const _ = CronJob.from({
            cronTime: '*/5 * * * * *',
            onTick: async () => await this.checkAccess(),
            start: true,
            waitForCompletion: true,
        });

        console.log('[CronJobService] Running.');
    }

    private async checkAccess() {
        await this.checkSleepTime(this.checkAccess.name);
        const accesses = await this.accessService.getAccess();
        console.log(accesses.length);
        if (accesses.length === 0) return;

        await this.accessService.checkAccess(accesses);
    }

    private async checkSleepTime(name: string) {
        const secondsLeft = getSleepSeconds(EnvConfig.timeZone);
        if (secondsLeft > 0) {
            const hours = (secondsLeft / 60 / 60).toFixed(1);
            console.log(`[CronJobService] Cron: ${name} paused for ${hours} hours.`);
        }

        await sleep(secondsLeft * 1000);
    }
}

export default CronJobService;

type TimeRange = [string, string];
type Scheduler = Record<number, TimeRange>;

const scheduler: Scheduler = {
    1: ['06:30', '23:00'],
    2: ['06:30', '23:00'],
    3: ['06:30', '23:00'],
    4: ['06:30', '23:00'],
    5: ['06:30', '23:00'],
    6: ['09:00', '20:00'],
    7: ['09:00', '14:00'],
};

export function getSleepSeconds(timeZone: string): number {
    const now = new Date();

    let chileOffset: number;

    switch (timeZone) {
        case 'GMT-3':
            chileOffset = -3;
            break;
        case 'GMT-4':
            chileOffset = -4;
            break;
        default:
            throw new Error('Unsupported timezone');
    }

    // Create a date in the Chile timezone using the timezone offset
    const utcTime = now.getTime();
    const chileTime = new Date(utcTime);
    // Set hours according to Chile's timezone
    chileTime.setUTCHours(now.getUTCHours() + chileOffset);

    const day = chileTime.getDay() === 0 ? 7 : chileTime.getDay();
    const [startStr, endStr] = scheduler[day];

    const [sh, sm] = startStr.split(':').map(Number);
    const [eh, em] = endStr.split(':').map(Number);

    const start = new Date(chileTime);
    start.setUTCHours(sh, sm, 0, 0);

    const end = new Date(chileTime);
    end.setUTCHours(eh, em, 0, 0);

    if (chileTime >= start && chileTime <= end) {
        return 0;
    } else if (chileTime < start) {
        return (start.getTime() - chileTime.getTime()) / 1000;
    } else {
        const nextDay = (day % 7) + 1;
        const [nextStartStr] = scheduler[nextDay];
        const [nsh, nsm] = nextStartStr.split(':').map(Number);

        const nextStart = new Date(chileTime);
        nextStart.setDate(nextStart.getDate() + 1);
        nextStart.setUTCHours(nsh, nsm, 0, 0);

        return (nextStart.getTime() - chileTime.getTime()) / 1000;
    }
}

export function sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

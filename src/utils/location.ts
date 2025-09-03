export function getLocation(location: string): string {
    switch (location) {
        case '102':
            return 'Espacio Urbano';
        case '104':
            return 'Calama';
        case '105':
            return 'Pacífico';
        case '106':
            return 'Arauco';
        case '107':
            return 'Iquique';
        case '108':
            return 'Angamos';
        default:
            return 'No location found';
    }
}

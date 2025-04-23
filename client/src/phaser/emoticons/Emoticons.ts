

export class Emoticons {
    public static getTextureAndFrame(emoticonStr: string): { texture: string; frame: string | number } {
        switch (emoticonStr) {
            case "attack": return { texture: 'actionicons', frame: 0 };
            case "forage": return { texture: 'actionicons', frame: 1 };
            case "chop": return { texture: 'actionicons', frame: 2 };
            case "mine": return { texture: 'actionicons', frame: 3 };
            case "harvest": return { texture: 'actionicons', frame: 4 };
            case "flee": return { texture: 'actionicons', frame: 10 };
            case "roam": return { texture: 'actionicons', frame: 12 };
            case "sell": return { texture: 'actionicons', frame: 8 };
            case "buy": return { texture: 'actionicons', frame: 9 };
            case "rest": return { texture: 'actionicons', frame: 11 };
            case "maintain": return { texture: 'actionicons', frame: 7 };
            case "rebuild": return { texture: 'actionicons', frame: 6 };
            case "dead": return { texture: 'emoticons', frame: 18 };
            default: {
                console.log(`No emoticon for '${emoticonStr}'`);
                return { texture: 'icons', frame: 'default' };
            }
        }
    }
}
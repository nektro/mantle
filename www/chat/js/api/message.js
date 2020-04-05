"use strict";
//

//
/** @type {Map<string,Map<string,Message} */
export const cache = new Map();

//
export class Message {
    constructor(o, ch_uid) {
        if (o === null) {
            this.is_null = true;
            return;
        }
        Object.assign(this, o);
        this.is_null = false;
        this.time = new moment.utc(this.time);
        cache.get(ch_uid).set(this.uuid, this);
    }
}

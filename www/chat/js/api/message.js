"use strict";
//

//
export const cache = new Map();

//
export class Message {
    constructor(o, ch) {
        if (o === null) {
            this.is_null = true;
            return;
        }
        Object.assign(this, o);
        this.is_null = false;
        cache.get(ch.uuid).set(this.uuid, this);
    }
}

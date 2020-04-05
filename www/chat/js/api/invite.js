"use strict";
//

//
/** @type {Map<string,Invite} */
export const cache = new Map();

//
export class Invite {
    constructor(o) {
        if (o === null) {
            this.is_null = true;
            return;
        }
        Object.assign(this, o);
        this.is_null = false;
        cache.set(this.uuid, this);
    }
}

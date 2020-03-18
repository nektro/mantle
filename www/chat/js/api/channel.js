"use strict";
//
//
export class Channel {
    constructor(o) {
        if (o === null) {
            this.is_null = true;
            return;
        }
        Object.assign(this, o);
        this.is_null = false;
    }
}

"use strict";
//
import { cache as cacheG } from "./message.js";

//
export const cache = new Map();

//
export class Channel {
    constructor(o) {
        if (o === null) {
            this.is_null = true;
            return;
        }
        Object.assign(this, o);
        this.is_null = false;
        cache.set(this.uuid, this);
        cacheG.set(this.uuid, new Map());
    }
}

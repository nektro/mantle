"use strict";
//
import * as api from "./index.js";

//
export const cache = new Map();

//
export class User {
    constructor(o) {
        if (o === null) {
            this.is_null = true;
            return;
        }
        Object.assign(this, o);
        this.is_null = false;
        cache.set(this.uuid, this);
    }
    /** @returns {Promise<Role[]>} */
    getRoles() {
        return this.roles.map((v) => api.M.roles.get(v));
    }
}

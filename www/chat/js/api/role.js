"use strict";
//

//
/** @type {Map<string,Role} */
export const cache = new Map();

//
export class Role {
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

//
new Role({
    uuid:"o",
    name:"Owner",
    color:"",
});

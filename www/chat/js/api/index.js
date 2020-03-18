"use strict";
//
import { User } from "./user.js";
import { Channel } from "./channel.js";
import { Role } from "./role.js";
import { Message } from "./message.js";

//
export {
    User,
    Channel,
    Role,
    Message,
};

function fetchE(endpoint) {
    return fetch(`./../api${endpoint}`).then((x) => x.json()).then((x) => {
        if (!x.success) {
            return Promise.reject(new Error(x.message));
        }
        return x;
    });
}

//
export const M = {
    users: {
        get: (uid) => {
            return fetchE(`/users/${uid}`).then((x) => {
                return new User(x.message);
            });
        },
        me: () => {
            return fetchE("/users/@me").then((x) => {
                return {
                    user: new User(x.message.me),
                    perms: x.message.me,
                };
            });
        },
    },
};

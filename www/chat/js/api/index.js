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
        return x.message;
    });
}

//
export const M = {
    users: {
        get: (uid) => {
            return fetchE(`/users/${uid}`).then((x) => {
                return new User(x);
            });
        },
        me: () => {
            return fetchE("/users/@me").then((x) => {
                return {
                    user: new User(x.me),
                    perms: x.me,
                };
            });
        },
    },
    roles: {
        get: () => {
            return fetchE("/roles").then((x) => {
                return x.map((y) => {
                    return new Role(y);
                });
            });
        },
    },
    channels: {
        me: () => {
            return fetchE("/channels/@me").then((x) => {
                return x.map((y) => {
                    return new Channel(y);
                });
            });
        },
    },
};

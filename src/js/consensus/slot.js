'use strict';
const delegates = 20;
const interval = 3;// second

function get_time(time) {
    if (time === undefined) {
        time = (new Date()).getTime();
    }
    var base_time = new Date(1548988864492).getTime();
    return Math.floor((time - base_time) / 1000);
}

function get_slot_number(time_stamp) {
    time_stamp = get_time(time_stamp);
    return Math.floor(time_stamp / interval);
}

module.exports = {
    interval,
    delegates,
    get_time,
    get_slot_number
};

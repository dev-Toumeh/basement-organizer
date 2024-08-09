import http from 'k6/http';
import { sleep } from 'k6';
import { check } from "k6";

// define configuration
export const options = {
    // define thresholds
    thresholds: {
        http_req_failed: [{ threshold: "rate<0.01", abortOnFail: true }], // availability threshold for error rate
        http_req_duration: ["p(99)<1000"], // Latency threshold for percentile
    },
    // define scenarios
    scenarios: {
        breaking: {
            executor: "ramping-vus",
            stages: [
                //{ duration: "5s", target: 20 },
                //{ duration: "5s", target: 200 },
                //{ duration: "5s", target: 400 },
                //{ duration: "5s", target: 600 },
                //{ duration: "5s", target: 800 },
                { duration: "5s", target: 1000 },
                { duration: "5s", target: 2000 },
                //{ duration: "5s", target: 14000 },
                //....
            ],
        },
    },
};

export default function() {
    const res = http.get("http://172.25.139.157:8000");
    //const res = http.get("http://172.25.139.157:8090");
    //const res = http.get("http://172.25.139.157:8000/api/v1/read/item/cc5f6abb-0e66-4dc7-a7ae-eb15deec285b")
    // check that response is 200
    check(res, {
        "response code was 200": (res) => res.status == 200,
    });
    //sleep(1);
}



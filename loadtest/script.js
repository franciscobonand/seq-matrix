import { sleep } from "k6";
import http from "k6/http";

export const options = {
  vus: 50,
  duration: "30s",
};

const data = JSON.stringify({
  letters: ["DUHBHB", "DBBUHD", "UBUUHU", "BHBDHH", "HDDDUB", "UDBDUH"],
});

export default function () {
  http.post("http://localhost:9001/sequence", data);
  // http.get("http://localhost:9001/stats");
  sleep(1);
}

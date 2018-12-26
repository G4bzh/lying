export default {
  name: "ZoneEdit",
  template: `
    <div>
      Hello {{ rrname }}
      TTL is {{ rrttl }}
    </div>
  `,
  props: {
    rrname: String,
    rrttl: Number,
    rrclass: String,
    rrtype: String,
    rrdata: String
  },
  data: function() {
    return {}
  }

};

export default {
  name: "RrEdit",
  template: `
    <div>
      <li>{{ rrindex }}- {{ rrname }} {{ rrttl }} {{ rrclass }} {{ rrtype }} {{ rrdata }}</li>
    </div>
  `,
  props: {
    rrindex: Number,
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

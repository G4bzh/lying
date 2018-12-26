export default {
  name: "RrEdit",
  template: `
    <div>
      <md-card>

        <md-card-header>
          {{ rrindex }}
        </md-card-header>

        <md-card-content>
           <md-field><md-input v-model="rrname_"></md-input></md-field>
          {{ rrttl }} {{ rrclass }} {{ rrtype }} {{ rrdata }}
        </md-card-content>
      </md-card>
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
    return {
      rrname_: this.rrname
    }
  }

};

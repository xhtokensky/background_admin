webpackJsonp([7],{ZP0H:function(e,t,a){"use strict";a("lMgl"),a("D9hl"),a("3ntp");var r=a("aiqZ"),l={props:{formdata:Object},components:{quillEditor:r.quillEditor},data:()=>({editorOption:{},rules:{model1:[{required:!0,message:"请输入",trigger:"blur"}],model2:[{required:!0,message:"请输入",trigger:"blur"}],model3:[{required:!0,message:"请输入",trigger:"blur"}],model4:[{required:!0,message:"请输入",trigger:"blur"}],model5:[{required:!0,message:"请输入",trigger:"blur"}],model6:[{required:!0,message:"请输入",trigger:"blur"}],model7:[{required:!0,message:"请输入",trigger:"blur"}],model11:[{required:!0,message:"请输入",trigger:"blur"}],model12:[{required:!0,message:"请输入",trigger:"blur"}],model51:[{required:!0,message:"请输入",trigger:"blur"}],model52:[{required:!0,message:"请输入",trigger:"blur"}],model53:[{required:!0,message:"请输入",trigger:"blur"}],editorModel1:[{required:!0,message:"请输入",trigger:"blur"}]}}),computed:{uploadUrl:function(){return this.$base_url+"upload/file"}},methods:{cancel(){this.$router.push({path:"message"})},submitForm(e){let t=this;this.$refs[e].validate(e=>{if(!e)return console.log("error submit!!"),!1;t.$emit("configSave",t.formdata)})}}},o={render:function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",{staticClass:"view "},[a("div",{staticStyle:{padding:".4rem 0 0 1rem"}},[a("el-form",{ref:"formdata",staticClass:"demo-ruleForm",attrs:{model:e.formdata,rules:e.rules,"label-width":"100px"}},[a("div",{staticClass:"dafter_ctt_dspBlock_clearBoth"},[a("div",{staticClass:"dfl",staticStyle:{"margin-right":".5rem",width:"50%"}},[e.formdata.label1?a("el-form-item",{attrs:{label:e.formdata.label1,prop:"model1"}},[a("el-input",{model:{value:e.formdata.model1,callback:function(t){e.$set(e.formdata,"model1",t)},expression:"formdata.model1"}})],1):e._e(),e._v(" "),e.formdata.label2?a("el-form-item",{attrs:{label:e.formdata.label2,prop:"model2"}},[a("el-input",{model:{value:e.formdata.model2,callback:function(t){e.$set(e.formdata,"model2",t)},expression:"formdata.model2"}})],1):e._e(),e._v(" "),e.formdata.label3?a("el-form-item",{attrs:{label:e.formdata.label3,prop:"model3"}},[a("el-input",{model:{value:e.formdata.model3,callback:function(t){e.$set(e.formdata,"model3",t)},expression:"formdata.model3"}})],1):e._e(),e._v(" "),e.formdata.label4?a("el-form-item",{attrs:{label:e.formdata.label4,prop:"model4"}},[a("el-input",{model:{value:e.formdata.model4,callback:function(t){e.$set(e.formdata,"model4",t)},expression:"formdata.model4"}})],1):e._e(),e._v(" "),e.formdata.label5?a("el-form-item",{attrs:{label:e.formdata.label5,prop:"model5"}},[a("el-input",{model:{value:e.formdata.model5,callback:function(t){e.$set(e.formdata,"model5",t)},expression:"formdata.model5"}})],1):e._e()],1)]),e._v(" "),a("div",{staticClass:"vue_quill_edito"},[e.formdata.editorLabel1?a("el-form-item",{attrs:{label:e.formdata.editorLabel1,prop:"editorModel1"}},[a("quill-editor",{ref:"myTextEditor",attrs:{options:e.editorOption},model:{value:e.formdata.editorModel1,callback:function(t){e.$set(e.formdata,"editorModel1",t)},expression:"formdata.editorModel1"}})],1):e._e()],1),e._v(" "),a("div",[e.formdata.label11?a("el-form-item",{attrs:{label:e.formdata.label11,prop:"model11"}},[a("el-radio-group",{model:{value:e.formdata.model11,callback:function(t){e.$set(e.formdata,"model11",t)},expression:"formdata.model11"}},[a("el-radio",{attrs:{label:0}},[e._v("全部用户")]),e._v(" "),a("el-radio",{attrs:{label:1}},[e._v("个人用户")])],1)],1):e._e(),e._v(" "),e.formdata.label12?a("el-form-item",{attrs:{label:e.formdata.label12}},[a("el-input",{attrs:{placeholder:"全部用户无需填写"},model:{value:e.formdata.model12,callback:function(t){e.$set(e.formdata,"model12",t)},expression:"formdata.model12"}})],1):e._e()],1),e._v(" "),a("div",{staticClass:"dpdl_140"},[a("el-button",{attrs:{type:"primary"},on:{click:e.cancel}},[e._v("取消")]),e._v(" "),a("el-button",{attrs:{type:"primary"},on:{click:function(t){return e.submitForm("formdata")}}},[e._v("保存")])],1)])],1)])},staticRenderFns:[]};var d=a("C7Lr")(l,o,!1,function(e){a("a5Nl")},"data-v-64361349",null);t.a=d.exports},a5Nl:function(e,t){},cuMd:function(e,t){},db6i:function(e,t,a){"use strict";Object.defineProperty(t,"__esModule",{value:!0});var r=a("t/xR"),l=a("D2g7"),o=a("ZP0H"),d={components:{position:r.a,searchV:l.a,messageEduit:o.a},data:()=>({isEduit:0,isDetail:0,searchData:{title:"记录册列表"},formdata:{label1:"标题",model1:"",editorLabel1:"内容",editorModel1:"",label11:"发送群体",model11:0,label12:"手机号",model12:""}}),mounted(){let e=this;e.$route.params.title?(e.isEduit=1,e.formdata.model1=e.$route.params.title,e.formdata.editorModel1=e.$route.params.content,e.formdata.model11=e.$route.params.type,e.formdata.model12=e.$route.params.phone?e.$route.params.phone:"",e.messageId=e.$route.params.messageId):e.isEduit=0},methods:{configSave(e){let t=this,a=t.$api+"tokensky/message/edit",r={title:e.model1,content:e.editorModel1,type:parseInt(e.model11),messageId:t.messageId};r.phone=e.model12?e.model12:"",t.$axios.post(a,r).then(e=>{e.data&&0==e.data.code&&(this.$message({message:"编辑成功",type:"success"}),t.$router.push({path:"message"}))})}}},s={render:function(){var e=this.$createElement,t=this._self._c||e;return t("div",{staticClass:"vc uservc"},[t("position",{attrs:{title:"生产过程管理 > OTC配置"}}),this._v(" "),t("div",{staticClass:"ctt"},[t("div",[t("messageEduit",{attrs:{formdata:this.formdata},on:{configSave:this.configSave}})],1)])],1)},staticRenderFns:[]};var m=a("C7Lr")(d,s,!1,function(e){a("cuMd")},"data-v-092594c8",null);t.default=m.exports}});
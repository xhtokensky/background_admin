webpackJsonp([28],{OOMD:function(e,a,t){"use strict";Object.defineProperty(a,"__esModule",{value:!0});var r=t("t/xR"),l=t("D2g7"),o={props:{formdata:Object},components:{},data:()=>({editorOption:{},rules:{model1:[{required:!0,message:"请输入",trigger:"blur"}],model2:[{required:!0,message:"请输入",trigger:"blur"}],model3:[{required:!0,message:"请输入",trigger:"blur"}],model4:[{required:!0,message:"请输入",trigger:"blur"}],model5:[{required:!0,message:"请输入",trigger:"blur"}],model6:[{required:!0,message:"请输入",trigger:"blur"}],model7:[{required:!0,message:"请输入",trigger:"blur"}],model11:[{required:!0,message:"请输入",trigger:"blur"}],model12:[{required:!0,message:"请输入",trigger:"blur"}],model51:[{required:!0,message:"请输入",trigger:"blur"}],model52:[{required:!0,message:"请输入",trigger:"blur"}],model53:[{required:!0,message:"请输入",trigger:"blur"}],editorModel1:[{required:!0,message:"请输入",trigger:"blur"}]}}),computed:{uploadUrl:function(){return this.$base_url+"upload/file"}},methods:{cancel(){this.$router.push({path:"otcBlacklist"})},submitForm(e){let a=this;this.$refs[e].validate(e=>{if(!e)return console.log("error submit!!"),!1;a.$emit("configSave",a.formdata)})}}},d={render:function(){var e=this,a=e.$createElement,t=e._self._c||a;return t("div",{staticClass:"view "},[t("div",{staticStyle:{padding:".4rem 0 0 1rem"}},[t("el-form",{ref:"formdata",staticClass:"demo-ruleForm",attrs:{model:e.formdata,rules:e.rules,"label-width":"250px"}},[t("div",{staticClass:"dafter_ctt_dspBlock_clearBoth"},[t("div",{staticClass:"dfl",staticStyle:{"margin-right":".5rem",width:"50%"}},[e.formdata.label1?t("el-form-item",{attrs:{label:e.formdata.label1,prop:"model1"}},[t("el-input",{model:{value:e.formdata.model1,callback:function(a){e.$set(e.formdata,"model1",a)},expression:"formdata.model1"}})],1):e._e(),e._v(" "),e.formdata.label2?t("el-form-item",{attrs:{label:e.formdata.label2,prop:"model2"}},[t("el-input",{model:{value:e.formdata.model2,callback:function(a){e.$set(e.formdata,"model2",a)},expression:"formdata.model2"}})],1):e._e(),e._v(" "),e.formdata.label3?t("el-form-item",{attrs:{label:e.formdata.label3,prop:"model3"}},[t("el-input",{model:{value:e.formdata.model3,callback:function(a){e.$set(e.formdata,"model3",a)},expression:"formdata.model3"}})],1):e._e(),e._v(" "),e.formdata.label4?t("el-form-item",{attrs:{label:e.formdata.label4,prop:"model4"}},[t("el-input",{model:{value:e.formdata.model4,callback:function(a){e.$set(e.formdata,"model4",a)},expression:"formdata.model4"}})],1):e._e(),e._v(" "),e.formdata.label5?t("el-form-item",{attrs:{label:e.formdata.label5,prop:"model5"}},[t("el-input",{model:{value:e.formdata.model5,callback:function(a){e.$set(e.formdata,"model5",a)},expression:"formdata.model5"}})],1):e._e()],1)]),e._v(" "),t("div",[e.formdata.label11?t("el-form-item",{attrs:{label:e.formdata.label11,prop:"model11"}},[t("el-radio-group",{model:{value:e.formdata.model11,callback:function(a){e.$set(e.formdata,"model11",a)},expression:"formdata.model11"}},[t("el-radio",{attrs:{label:1}},[e._v("禁止登陆")]),e._v(" "),t("el-radio",{attrs:{label:2}},[e._v("禁止交易")])],1)],1):e._e()],1),e._v(" "),t("div",{staticClass:"dpdl_140"},[t("el-button",{attrs:{type:"primary"},on:{click:e.cancel}},[e._v("取消")]),e._v(" "),t("el-button",{attrs:{type:"primary"},on:{click:function(a){return e.submitForm("formdata")}}},[e._v("保存")])],1)])],1)])},staticRenderFns:[]};var s=t("C7Lr")(o,d,!1,function(e){t("xfK8")},"data-v-24611d16",null).exports,m={components:{position:r.a,searchV:l.a,blackEduForm:s},data:()=>({isEduit:0,isDetail:0,searchData:{title:"记录册列表"},formdata:{label1:"手机号",model1:"",label2:"封停期限(h)",model2:"",label11:"封停类型",model11:1,label12:"手机号",model12:""}}),mounted(){let e=this;e.$route.params&&e.$route.params.phone?(e.isEduit=1,e.id=e.$route.params.id,e.formdata.model1=e.$route.params.phone,e.formdata.model2=e.$route.params.periodTime/3600,e.formdata.model11=e.$route.params.balckType):e.isEduit=0},methods:{configSave(e){let a=this,t=a.$api+"role/blackList/edit",r={phone:e.model1,periodTime:3600*parseInt(e.model2),balckType:parseInt(e.model11)};(a.isEduit=1)?(r.id=a.id,a.$axios.post(t,r).then(e=>{e.data&&0==e.data.code&&(this.$message({message:e.data.msg,type:"success"}),a.$router.push({path:"otcBlacklist"}))})):a.$axios.post(t,r).then(e=>{e.data&&0==e.data.code&&(this.$message({message:e.data.msg,type:"success"}),a.$router.push({path:"otcBlacklist"}))})}}},i={render:function(){var e=this.$createElement,a=this._self._c||e;return a("div",{staticClass:"vc uservc"},[a("position",{attrs:{title:"生产过程管理 > OTC配置"}}),this._v(" "),a("div",{staticClass:"ctt"},[a("div",[a("blackEduForm",{attrs:{formdata:this.formdata},on:{configSave:this.configSave}})],1)])],1)},staticRenderFns:[]};var u=t("C7Lr")(m,i,!1,function(e){t("dv4P")},"data-v-18967836",null);a.default=u.exports},dv4P:function(e,a){},xfK8:function(e,a){}});
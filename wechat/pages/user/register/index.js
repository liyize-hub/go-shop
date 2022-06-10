import { request, getUserProfile, showToast } from "../../../request/index.js";
var config = require('../../../config/config.js');
Page({

  /**
   * 页面的初始数据
   */
  data: {
    mobile: 0,
    address: "",
    token: ""
  },
  UserInfo: {
    phone: 0,
    address: "",
    img: "",
    name: "",
    token: ""
  },

  onLoad(options) {
    //从登录页面获取token
    console.log(options)
    this.setData({
      token: options.token
    })
  },

  //获取手机号
  mobileblur(e) {
    var content = e.detail.value;
    this.setData({
      mobile: content
    });
  },

  //获取收货地址
  addressblur(e) {
    var content = e.detail.value;
    this.setData({
      address: content
    })
  },

  async register(e) {
    if (this.data.mobile == 0) {
      await showToast({ title: '手机号为空' });
      return
    }
    if (this.data.address == "") {
      await showToast({ title: '收货地址为空' });
      return
    }
    const res = await getUserProfile({ desc: "用于完善会员资料" });
    this.UserInfo.phone = this.data.mobile;
    this.UserInfo.address = this.data.address;
    this.UserInfo.token = this.data.token;
    this.UserInfo.img = res.userInfo.avatarUrl;
    this.UserInfo.name = res.userInfo.nickName;
    const resp = await request({
      url: config.api.weAppRegister, data: this.UserInfo, method: "POST", header:
        { "Content-Type": "application/json;charset=UTF-8" }
    });

    if (resp.errNo == 0) {
      await showToast({ title: '注册成功' });
    } else {
      await showToast({ title: '用户注册失败，请重试' });
    }

    wx.navigateBack({
      delta: 1
    });
  }
})

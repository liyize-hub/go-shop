// pages/login/index.js
import { getUserProfile, request } from "../../../request/index.js";
var config = require('../../../config/config.js');
Page({

  data: {
    userInfo: [],
    mustgetUserInfo: true,
    isregister: false,
  },

  /**
     * 生命周期函数--监听页面加载
     */
  onLoad(options) {
  },

  //检查用户是否已注册
  check() {
    wx.login({
      success: async (res) => {
        if (res.code) {
          //发起网络请求
          const resp = await request({ url: config.api.weAppLogin, data: { code: res.code } });
          if (resp.errNo === 2) {
            // 3 跳转到 注册页面
            wx.navigateTo({
              url: '/pages/user/register/index?token=' + resp.data
            });
          }
          if (resp.errNo === 0) {
            wx.setStorageSync("userInfo", { data: resp.data });
            wx.navigateBack({
              delta: 1
            });
          }
        } else {
          console.log('登录失败！' + res.errMsg)
        }
      }
    })
  },
})
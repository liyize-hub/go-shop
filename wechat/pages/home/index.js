import { request } from "../../request/index.js";
// pages/home/index.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    swiperList: [
      {
        id: '1',
        image: `https://go-shop-img-1311006175.cos.ap-beijing.myqcloud.com/%E8%A1%8C%E6%9D%8E%E7%AE%B1-%E8%BD%AE%E6%92%AD%E5%9B%BE.png`,
      },
      {
        id: '2',
        image: `https://go-shop-img-1311006175.cos.ap-beijing.myqcloud.com/%E8%80%B3%E6%9C%BA-%E8%BD%AE%E6%92%AD%E5%9B%BE.png`,
      },
      {
        id: '3',
        image: `https://go-shop-img-1311006175.cos.ap-beijing.myqcloud.com/%E6%AF%9B%E6%AF%AF-%E8%BD%AE%E6%92%AD%E5%9B%BE.png`,
      },
      {
        id: '4',
        image: `https://go-shop-img-1311006175.cos.ap-beijing.myqcloud.com/%E6%9E%81%E5%85%89%E7%9B%92%E5%AD%90-%E8%BD%AE%E6%92%AD%E5%9B%BE.png`,
      },
    ],
    productList: []
  },


  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady() {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow() {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide() {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload() {

  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh() {

  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom() {

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage() {

  }
})
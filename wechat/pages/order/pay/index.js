import { showModal } from "../../../request/index.js";

Page({
    data: {
        cart: [],
        totalPrice: 0,
        totalNum: 0,
        scrollTop: 0,
        minutes: '00', //分
        seconds: '00', //秒
        //leftTime: 10 * 60 //倒计时10分钟
    },
    onLoad: function () {
        this.countTime();
    },
    onShow() {
        // 1 获取缓存中的购物车数据
        let cart = wx.getStorageSync("cart") || [];
        // 过滤后的购物车数组
        cart = cart.filter(v => v.checked);
        // 1 总价格 总数量
        let totalPrice = 0;
        let totalNum = 0;
        cart.forEach(v => {
            totalPrice += v.num * v.price;
            totalNum += v.num;
        })
        this.setData({
            cart,
            totalPrice, totalNum
        });
    },

    async countTime() {
        let minutes, seconds;
        let that = this;
        const app = getApp();
        if (app.globalData.leftTime >= 0) {
            minutes = Math.floor(app.globalData.leftTime / 60 % 60);
            seconds = Math.floor(app.globalData.leftTime % 60);
            seconds = seconds < 10 ? "0" + seconds : seconds
            minutes = minutes < 10 ? "0" + minutes : minutes
            that.setData({
                minutes,
                seconds
            })
            app.globalData.leftTime--;
            setTimeout(that.countTime, 1000);
        } else {
            await showModal({ content: "订单已过期" });
            // 清空购物车中过期的订单
            let cartLeft = that.data.cart.filter(v => !v.checked);
            wx.setStorageSync("cart", cartLeft);
        }
    },

})
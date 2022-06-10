import { showModal, showToast, request } from "../../request/index.js";
var config = require('../../config/config.js');
Page({
    data: {
        cart: [],
        allChecked: false,
        totalPrice: 0,
        totalNum: 0,
        scrollTop: 0
    },
    onLoad: function () {
    },
    onShow() {
        // 1 获取缓存中的购物车数据
        const cart = wx.getStorageSync("cart") || [];
        this.setCart(cart);
    },

    // 商品的选中
    handeItemChange(e) {
        // 1 获取被修改的商品的id
        const goods_id = e.currentTarget.dataset.id;
        // 2 获取购物车数组 
        let { cart } = this.data;
        // 3 找到被修改的商品对象
        let index = cart.findIndex(v => v.id === goods_id);
        // 4 选中状态取反
        cart[index].checked = !cart[index].checked;

        this.setCart(cart);
    },

    // 商品全选 反选功能
    handleItemAllCheck() {
        // 1 获取data中的数据
        let { cart, allChecked } = this.data;
        // 2 修改值
        allChecked = !allChecked;
        // 3 循环修改cart数组 中的商品选中状态
        cart.forEach(v => v.checked = allChecked);
        // 4 把修改后的值 填充回data或者缓存中
        this.setCart(cart);
    },

    // 商品数量的编辑功能
    async handleItemNumEdit(e) {
        // 1 获取传递过来的参数 
        const { operation, id } = e.currentTarget.dataset;
        // 2 获取购物车数组
        let { cart } = this.data;
        // 3 找到需要修改的商品的索引
        const index = cart.findIndex(v => v.id === id);
        // 4 判断是否要执行删除
        if (cart[index].num === 1 && operation === -1) {
            // 4.1 弹窗提示
            const res = await showModal({ content: "您是否要删除？" });
            if (res.confirm) {
                cart.splice(index, 1);
                this.setCart(cart);
            }
        } else {
            // 4  进行修改数量
            cart[index].num += operation;
            // 5 设置回缓存和data中
            this.setCart(cart);
        }
    },

    // 设置购物车状态同时 重新计算 底部工具栏的数据 全选 总价格 购买的数量
    setCart(cart) {
        let allChecked = true;
        // 1 总价格 总数量
        let totalPrice = 0;
        let totalNum = 0;
        cart.forEach(v => {
            if (v.checked) {
                totalPrice += v.num * v.price;
                totalNum += v.num;
            } else {
                allChecked = false;
            }
        })
        // 判断数组是否为空
        allChecked = cart.length != 0 ? allChecked : false;
        this.setData({
            cart,
            totalPrice, totalNum, allChecked
        });
        wx.setStorageSync("cart", cart);
    },

    // 点击 结算 
    async handlePay() {
        // 1 判断收货地址
        const { totalNum } = this.data;
        // 2 判断用户有没有选购商品
        if (totalNum === 0) {
            await showToast({ title: "您还没有选购商品" });
            return;
        }
        // 3 判断本地缓存有没有userInfo数据
        const userInfo = wx.getStorageSync("userInfo");
        if (!userInfo) {
            wx.navigateTo({
                url: '/pages/user/login/index'
            });
        }
        // 4 判断用户是否登录过期
        const resp = await request({ url: config.api.weAppCheck, data: { token: userInfo.data.token, id: userInfo.data.uid } });
        if (resp.errNo != 0) {
            wx.navigateTo({
                url: '/pages/user/login/index'
            });
        }
        // 5 跳转到 支付页面
        wx.navigateTo({
            url: '/pages/order/pay/index'
        });

    }

})
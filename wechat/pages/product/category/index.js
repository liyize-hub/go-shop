
import { request } from "../../../request/index.js";
var config = require('../../../config/config.js');
// pages/product/category/index.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    categoryList: [],
    productList: [],
    // 被点击的左侧的菜单
    currentIndex: -1,
    // 右侧内容的滚动条距离顶部的距离
    scrollTop: 0
  },
  // 总页数
  totalPages: 1,
  // 接口要的参数
  QueryParams: {
    id: 0,
    page: 1,
    perpage: 10
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {

    this.isStorage("category")
  },

  // 1 使用es7的async await来发送请求 获取商品种类列表
  async getCategoryList() {
    const res = await request({ url: config.api.GetCategoryAll });
    // 把接口的数据存入到本地存储中
    wx.setStorageSync("category", { time: Date.now(), data: res.data.items });

    this.setData({
      categoryList: res.data.items
    })
  },

  // 2 获取商品列表
  async getProductList() {
    const res = await request({ url: config.api.GetCategoryAll, data: this.QueryParams });
    if (res.data.total === 0) {
      return
    }
    // 获取 总条数
    const total = res.data.total;
    // 计算总页数
    this.totalPages = Math.ceil(total / this.QueryParams.perpage);

    this.setData({
      productList: [...this.data.productList, ...res.data.items]
    })

    wx.setStorageSync("商品种类" + this.QueryParams.id.toString(), { time: Date.now(), data: this.data.productList });
  },

  // 左侧菜单的点击事件
  async handleItemTap(e) {
    /* 
    1 获取被点击的标题身上的索引
    2 给data中的currentIndex赋值就可以了
    3 根据不同的索引来渲染右侧的商品内容
     */
    if (this.QueryParams.id == e.currentTarget.dataset.id) {
      return
    }
    //将商品列表和页数清零
    this.setData({
      productList: [],
    })
    this.QueryParams.page = 1;

    this.QueryParams.id = e.currentTarget.dataset.id;
    this.isStorage("商品种类" + this.QueryParams.id.toString())

    this.setData({
      currentIndex: (this.QueryParams.id - 1),
      // 重新设置 右侧内容的scroll-view标签的距离顶部的距离
      scrollTop: 0
    })
  },

  // 页面上滑 滚动条触底事件
  onReachBottom() {
    //  1 判断还有没有下一页数据
    if (this.QueryParams.page >= this.totalPages) {
      // 没有下一页数据
      wx.showToast({ title: '没有下一页数据' });

    } else {
      // 还有下一页数据
      this.QueryParams.page++;
      this.getProductList();
    }
  },

  // 下拉刷新事件 
  onPullDownRefresh() {
    // 1 重置数组
    this.setData({
      productList: [],
      currentIndex: -1,
      scrollTop: 0
    })
    // 2 重置页码
    this.QueryParams.page = 1;
    this.QueryParams.id = 0;
  },

  isStorage(key) {
    //  1 获取本地存储中的数据  (小程序中也是存在本地存储 技术)
    const Storage = wx.getStorageSync(key);

    // 2 判断
    if (!Storage) {
      // 不存在  发送请求获取数据
      if (key == "category") {
        this.getCategoryList();
      } else {
        this.getProductList();
      }
    } else {
      // 有旧的数据 定义过期时间  10s 改成 2分钟
      if (Date.now() - Storage.time > 1000 * 60 * 2) {
        // 清除缓存记录
        wx.removeStorageSync(key);
        // 重新发送请求
        if (key == "category") {
          this.getCategoryList();
        } else {
          this.getProductList();
        }
      } else {
        // 可以使用旧数据
        if (key == "category") {
          this.setData({
            categoryList: Storage.data
          })
        } else {
          this.setData({
            productList: Storage.data
          })
        }
      }
    }

  },

  // 点击 加入购物车
  handleCartAdd(e) {
    let product = e.currentTarget.dataset
    // 1 获取缓存中的购物车 数组
    let cart = wx.getStorageSync("cart") || [];
    // 2 判断 商品对象是否存在于购物车数组中
    let index = cart.findIndex(v => v.id === product.id);
    if (index === -1) {
      //3  不存在 第一次添加
      product = e.currentTarget.dataset
      product.num = 1;
      product.checked = true;
      cart.push(product);
    } else {
      // 4 已经存在购物车数据 执行 num++
      cart[index].num++;
    }
    // 5 把购物车重新添加回缓存中
    wx.setStorageSync("cart", cart);
    // 6 弹窗提示
    wx.showToast({
      title: '加入成功',
      icon: 'success',
      duration: 200,
      mask: true
    });
  },

})

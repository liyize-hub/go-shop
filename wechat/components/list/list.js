// components/list/list.js
Component({
  /**
   * 组件的属性列表
   */
  properties: {
    productList: {
      type: Array,
      value: []
    }
  },
  methods: {

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
    }

  }


})

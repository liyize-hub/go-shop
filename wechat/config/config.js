var url = 'http://localhost:8080';

var config = {
    name: "本地购物商城",
    wemallSession: "wemallSession",
    static: {
        imageDomain: url
    },
    api: {
        weAppLogin: '/user/login',
        weAppRegister: '/user/register',
        weAppCheck: '/user/check',
        setWeAppUser: '/setWeAppUser',
        GetCategoryAll: '/category',
        PostProductList: '/product/select',
        addToCart: '/cart/create'
    }
};

for (var key in config.api) {
    config.api[key] = url + config.api[key];
}

module.exports = config;
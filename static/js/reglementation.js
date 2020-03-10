$( document ).ready(function() {
    
    /* CHANGE LES COULEURS DU BLOC JURIDIQUE AU CLICK*/
    $(".itemJuridique").click(function() {
        $(".itemJuridique").toggleClass("itemJuridiqueActived");
        $(".lienJuridique").toggleClass("lienJuridiqueActived");
    });

    /* CHANGE LES COULEURS DES ITEMS AU CLICK*/
    $(".itemsFacultatif").click(function(){
        $(this).toggleClass("itemsFacultatifActived");
        $(this).find('.stroke').toggleClass("strokeActived");
        $(this).find('.fill').toggleClass("fillActived");
    });

});
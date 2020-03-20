$( document ).ready(function() {
	var swiper = new Swiper('.swiper-container', {
		slidesPerView: 2,
		loop: false,
		spaceBetween: 0,
		breakpoints: {
			319: {
			slidesPerView: 2,
			spaceBetween: 1,
			},
			374: {
			slidesPerView: 4,
			spaceBetween: 12,
			},
			413: {
			slidesPerView: 2,
			spaceBetween: 20,
			},
		}
	});	

	// OUVRIR SIDE MENU
	$('.menuHeader').click(function(){
		$('.aside').toggle();
		$('header').toggleClass("flou");
		$('section').toggleClass("flou");
		$('footer').toggleClass("flou");
	});
	// FERMER SIDE MENU
	$('.closeBtn').click(function(){
		$('.aside').toggle();
		$('header').toggleClass("flou");
		$('section').toggleClass("flou");
		$('footer').toggleClass("flou");
	});

    /* CHANGE LES COULEURS DES ITEMS AU CLICK*/
    $(".swiper-slide").click(function(){
        $(this).toggleClass("itemActived");
		$(this).find('.icone').toggleClass("iconeActived");

		// afficher l'ordre des criteres achats
		var n = $(".itemActived").length;
		$(this).parent().find(".numeroItem").text(n);
		$(this).not(".itemActived").parent().find(".numeroItem").text('');
		
		// cganger couleur svg
		$(this).find('.stroke').toggleClass("strokeActived");
        $(this).find('.fill').toggleClass("fillActived");
	});

});
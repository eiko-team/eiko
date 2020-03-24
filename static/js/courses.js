$( document ).ready(function() {
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
	// OUVRIR ET FERMER LE CHEVRON D'UN PRODUIT
	$('.dropdownProduit').click(function(){
		$(this).toggleClass("chevronReserve");
	})

});
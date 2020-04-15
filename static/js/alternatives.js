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
	$('.dropdownHeader2').click(function(){
		$(this).toggleClass("chevronReserve");
	})
});
	// AFFICHER L'ICONE DE LA PAGE EN COURS EN ROUGE
	function activeFooter() {
		var title = $('form').hasClass("searchBarNav");

		if (title = true) {
			$('.footerAlternative').css("fill","#B53D00");
			$('.alternativeLien').css("color","#B53D00");
		}
		else {
			return false;
		}
	};
	activeFooter();	
		

	
